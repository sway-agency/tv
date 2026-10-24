package fr.sway.accueil;

import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.content.pm.ResolveInfo;
import android.graphics.drawable.BitmapDrawable;
import android.graphics.drawable.Drawable;

import java.text.Collator;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

final class AppCatalog {
    private static final String HIDDEN_PACKAGE = "com.tcl.partnercustomizer";
    private static final String FIREPLACE_BANNER = "fireplace/banner.jpg";

    private final Context context;
    private List<LauncherApp> cached;
    private int version;

    AppCatalog(Context context) {
        this.context = context;
    }

    void invalidate() {
        cached = null;
        version++;
    }

    int version() {
        return version;
    }

    List<LauncherApp> tvApps() {
        if (cached == null) {
            cached = Collections.unmodifiableList(load());
        }
        return cached;
    }

    List<LauncherApp> pick(List<String> packageNames) {
        List<LauncherApp> picked = new ArrayList<>();
        for (String packageName : packageNames) {
            for (LauncherApp app : tvApps()) {
                if (app.packageName.equals(packageName)) {
                    picked.add(app);
                    break;
                }
            }
        }
        return picked;
    }

    private List<LauncherApp> load() {
        PackageManager packageManager = context.getPackageManager();
        Intent query = new Intent(Intent.ACTION_MAIN).addCategory(Intent.CATEGORY_LEANBACK_LAUNCHER);
        List<LauncherApp> apps = new ArrayList<>();
        for (ResolveInfo resolveInfo : packageManager.queryIntentActivities(query, 0)) {
            ActivityInfo activity = resolveInfo.activityInfo;
            if (activity.packageName.equals(context.getPackageName()) || activity.packageName.equals(HIDDEN_PACKAGE)) {
                continue;
            }
            Drawable banner = activity.loadBanner(packageManager);
            Drawable icon = banner == null ? resolveInfo.loadIcon(packageManager) : null;
            apps.add(new LauncherApp(activity.packageName, resolveInfo.loadLabel(packageManager).toString(),
                    banner, icon, launchIntentFor(activity)));
        }
        apps.add(fireplace());
        Collator collator = Collator.getInstance();
        apps.sort((left, right) -> collator.compare(left.label, right.label));
        return apps;
    }

    private LauncherApp fireplace() {
        Drawable banner = new BitmapDrawable(context.getResources(), AssetImages.opaque(context, FIREPLACE_BANNER));
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
