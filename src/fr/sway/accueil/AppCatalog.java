package fr.sway.accueil;

import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.content.pm.ResolveInfo;
import android.content.res.Resources;
import android.graphics.drawable.Drawable;
import android.util.DisplayMetrics;

import java.io.IOException;
import java.io.InputStream;
import java.text.Collator;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

final class AppCatalog {
    private static final List<String> HIDDEN_PACKAGES = Arrays.asList("com.tcl.partnercustomizer");
    private static final String FIREPLACE_BANNER = "fireplace/banner.jpg";

    private final Context context;

    AppCatalog(Context context) {
        this.context = context;
    }

    List<LauncherApp> tvApps() {
        PackageManager packageManager = context.getPackageManager();
        Intent query = new Intent(Intent.ACTION_MAIN).addCategory(Intent.CATEGORY_LEANBACK_LAUNCHER);
        List<LauncherApp> apps = new ArrayList<>();
        for (ResolveInfo resolveInfo : packageManager.queryIntentActivities(query, 0)) {
            ActivityInfo activity = resolveInfo.activityInfo;
            if (activity.packageName.equals(context.getPackageName()) || HIDDEN_PACKAGES.contains(activity.packageName)) {
                continue;
            }
            apps.add(new LauncherApp(
                    activity.packageName,
                    resolveInfo.loadLabel(packageManager).toString(),
                    sharpBanner(packageManager, activity),
                    resolveInfo.loadIcon(packageManager),
                    launchIntentFor(activity)));
        }
        apps.add(fireplace());
        Collator collator = Collator.getInstance();
        apps.sort((left, right) -> collator.compare(left.label, right.label));
        return apps;
    }

    List<LauncherApp> pick(List<LauncherApp> apps, List<String> packageNames) {
        List<LauncherApp> picked = new ArrayList<>();
        for (String packageName : packageNames) {
            for (LauncherApp app : apps) {
                if (app.packageName.equals(packageName)) {
                    picked.add(app);
                    break;
                }
            }
        }
        return picked;
    }

    private static Drawable sharpBanner(PackageManager packageManager, ActivityInfo activity) {
        int bannerResource = activity.getBannerResource();
        if (bannerResource != 0) {
            try {
                Resources resources = packageManager.getResourcesForApplication(activity.applicationInfo);
                return resources.getDrawableForDensity(bannerResource, DisplayMetrics.DENSITY_XXXHIGH, null);
            } catch (PackageManager.NameNotFoundException | Resources.NotFoundException ignored) {
            }
        }
        return activity.loadBanner(packageManager);
    }

    private LauncherApp fireplace() {
        Drawable banner = null;
        try (InputStream stream = context.getAssets().open(FIREPLACE_BANNER)) {
            banner = Drawable.createFromStream(stream, FIREPLACE_BANNER);
        } catch (IOException ignored) {
        }
        Intent intent = new Intent(context, FireplaceActivity.class).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK);
        return new LauncherApp(context.getPackageName(), "Cheminée", banner, null, intent);
    }

    private static Intent launchIntentFor(ActivityInfo activity) {
        return new Intent(Intent.ACTION_MAIN)
                .addCategory(Intent.CATEGORY_LEANBACK_LAUNCHER)
                .setComponent(new ComponentName(activity.packageName, activity.name))
                .addFlags(Intent.FLAG_ACTIVITY_NEW_TASK | Intent.FLAG_ACTIVITY_RESET_TASK_IF_NEEDED);
    }
}
