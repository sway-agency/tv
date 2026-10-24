package fr.sway.accueil;

import android.content.Intent;
import android.graphics.drawable.Drawable;

final class LauncherApp {
    final String packageName;
    final String label;
    final Drawable banner;
    final Drawable icon;
    final Intent launchIntent;

    LauncherApp(String packageName, String label, Drawable banner, Drawable icon, Intent launchIntent) {
        this.packageName = packageName;
        this.label = label;
        this.banner = banner;
        this.icon = icon;
        this.launchIntent = launchIntent;
    }
}
