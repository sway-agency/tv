package fr.sway.accueil;

import android.app.Activity;
import android.content.ActivityNotFoundException;
import android.content.Intent;
import android.graphics.Color;
import android.graphics.drawable.GradientDrawable;
import android.os.Bundle;
import android.provider.Settings;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.FrameLayout;
import android.widget.ImageView;
import android.widget.LinearLayout;

import java.util.ArrayList;
import java.util.List;

public class HomeActivity extends Activity {
    private static final String GRID_ASSET = "branding/grid.png";

    private AppCatalog catalog;
    private DockConfig dockConfig;
    private TileFactory tiles;
    private WallpaperRotator wallpaperRotator;
    private SnowView snow;
    private LinearLayout dock;
    private AppDrawer drawer;
    private int screenHeight;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        int screenWidth = getResources().getDisplayMetrics().widthPixels;
        screenHeight = getResources().getDisplayMetrics().heightPixels;
        catalog = new AppCatalog(this);
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
        dock = dock();
        root.addView(dock, dockParams());
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
    protected void onPause() {
        super.onPause();
        wallpaperRotator.stop();
        snow.stop();
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        closeDrawer();
        focusFirstDockTile();
    }

    @Override
    public void onBackPressed() {
        closeDrawer();
    }

    private void populateDock() {
        dock.removeAllViews();
        for (LauncherApp app : catalog.pick(catalog.tvApps(), dockConfig.packages())) {
            dock.addView(tiles.appTile(app, () -> launch(app.launchIntent)));
        }
        dock.addView(tiles.iconTile(AssetImages.bitmap(this, GRID_ASSET), "Applis", "Toutes les applis", this::openDrawer));
        if (!drawer.isOpen()) {
            focusFirstDockTile();
        }
    }

    private void openDrawer() {
        List<View> appTiles = new ArrayList<>();
        for (LauncherApp app : catalog.tvApps()) {
            appTiles.add(tiles.appTile(app, () -> launch(app.launchIntent)));
        }
        dock.setDescendantFocusability(ViewGroup.FOCUS_BLOCK_DESCENDANTS);
        drawer.open(appTiles);
    }

    private void closeDrawer() {
        if (drawer == null || !drawer.isOpen()) {
            return;
        }
        drawer.close();
        dock.setDescendantFocusability(ViewGroup.FOCUS_AFTER_DESCENDANTS);
        focusFirstDockTile();
    }

    private void focusFirstDockTile() {
        if (dock.getChildCount() > 0) {
            dock.getChildAt(0).requestFocus();
        }
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

    private LinearLayout dock() {
        LinearLayout row = new LinearLayout(this);
        row.setOrientation(LinearLayout.HORIZONTAL);
        row.setGravity(Gravity.CENTER);
        row.setClipChildren(false);
        row.setClipToPadding(false);
        return row;
    }

    private FrameLayout.LayoutParams dockParams() {
        FrameLayout.LayoutParams params = new FrameLayout.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT, Gravity.BOTTOM);
        params.bottomMargin = screenHeight / 14;
        return params;
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
