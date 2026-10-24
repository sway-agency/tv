package fr.sway.accueil;

import android.app.Activity;
import android.content.ActivityNotFoundException;
import android.content.Intent;
import android.graphics.Color;
import android.graphics.drawable.GradientDrawable;
import android.os.Bundle;
import android.provider.Settings;
import android.view.Gravity;
import android.view.KeyEvent;
import android.view.View;
import android.view.ViewGroup;
import android.widget.FrameLayout;
import android.widget.ImageView;

import java.util.ArrayList;
import java.util.List;

public class HomeActivity extends Activity {
    private AppCatalog catalog;
    private DockConfig dockConfig;
    private TileFactory tiles;
    private WallpaperRotator wallpaperRotator;
    private SnowView snow;
    private Dock dock;
    private AppDrawer drawer;
    private PackageChanges packageChanges;
    private List<View> drawerTiles = new ArrayList<>();
    private int drawerVersion = -1;
    private int screenHeight;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        int screenWidth = getResources().getDisplayMetrics().widthPixels;
        screenHeight = getResources().getDisplayMetrics().heightPixels;
        catalog = new AppCatalog(this);
        packageChanges = new PackageChanges(catalog::invalidate);
        packageChanges.register(this);
        dockConfig = new DockConfig(this);
        tiles = new TileFactory(this, screenWidth * 12 / 100);

        FrameLayout root = new FrameLayout(this);
        root.setBackgroundColor(Color.BLACK);
        root.setClipChildren(false);
        ImageView backLayer = wallpaperLayer();
        ImageView frontLayer = wallpaperLayer();
        root.addView(backLayer, matchParent());
        root.addView(frontLayer, matchParent());
        root.addView(shade(), matchParent());
        snow = new SnowView(this);
        root.addView(snow, matchParent());

        TopBar topBar = new TopBar(this, screenWidth, screenHeight);
        root.addView(topBar.clock(), wrapContent(Gravity.TOP | Gravity.START));
        root.addView(topBar.corner(this::openSettings), wrapContent(Gravity.TOP | Gravity.END));
        dock = new Dock(this, tiles, this::launch, this::openDrawer);
        root.addView(dock, Dock.layoutParams(screenHeight));
        drawer = new AppDrawer(this, screenHeight);
        root.addView(drawer, matchParent());

        setContentView(root);
        wallpaperRotator = new WallpaperRotator(getAssets(), backLayer, frontLayer);
    }

    @Override
    protected void onResume() {
        super.onResume();
        populateDock();
        Motion.riseIn(dock, screenHeight / 12f);
        wallpaperRotator.start();
        snow.start();
    }

    @Override
    protected void onDestroy() {
        unregisterReceiver(packageChanges);
        super.onDestroy();
    }

    @Override
    protected void onPause() {
        super.onPause();
        wallpaperRotator.stop();
        snow.stop();
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        closeDrawer();
        dock.focusFirst();
    }

    @Override
    public boolean dispatchKeyEvent(KeyEvent event) {
        if (event.getKeyCode() != KeyEvent.KEYCODE_BACK) {
            return super.dispatchKeyEvent(event);
        }
        if (event.getAction() == KeyEvent.ACTION_UP) {
            closeDrawer();
        }
        return true;
    }

    private void populateDock() {
        List<String> packages = dockConfig.packages();
        boolean rebuilt = dock.show(catalog.pick(packages), catalog.version() + ":" + packages);
        if (rebuilt && !drawer.isOpen()) {
            dock.focusFirst();
        }
    }

    private void openDrawer() {
        if (drawerVersion != catalog.version()) {
            drawerVersion = catalog.version();
            drawerTiles = new ArrayList<>();
            for (LauncherApp app : catalog.tvApps()) {
                drawerTiles.add(tiles.appTile(app, () -> launch(app.launchIntent)));
            }
        }
        dock.setBlocked(true);
        drawer.open(drawerTiles);
    }

    private void closeDrawer() {
        if (drawer == null || !drawer.isOpen()) {
            return;
        }
        drawer.close();
        dock.setBlocked(false);
        dock.focusFirst();
    }

    private void openSettings() {
        launch(new Intent(Settings.ACTION_SETTINGS).addFlags(Intent.FLAG_ACTIVITY_NEW_TASK));
    }

    private void launch(Intent intent) {
        try {
            startActivity(intent);
        } catch (ActivityNotFoundException | SecurityException ignored) {
        }
    }

    private ImageView wallpaperLayer() {
        ImageView layer = new ImageView(this);
        layer.setScaleType(ImageView.ScaleType.CENTER_CROP);
        return layer;
    }

    private View shade() {
        View shade = new View(this);
        shade.setBackground(new GradientDrawable(GradientDrawable.Orientation.TOP_BOTTOM,
                new int[] {Color.argb(70, 0, 0, 0), Color.TRANSPARENT, Color.argb(170, 0, 0, 0)}));
        return shade;
    }

    private static FrameLayout.LayoutParams wrapContent(int gravity) {
        return new FrameLayout.LayoutParams(ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT, gravity);
    }

    private static FrameLayout.LayoutParams matchParent() {
        return new FrameLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT);
    }
}
