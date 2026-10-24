package fr.sway.accueil;

import android.content.Context;
import android.content.Intent;
import android.view.Gravity;
import android.view.ViewGroup;
import android.widget.FrameLayout;
import android.widget.LinearLayout;

import java.util.List;
import java.util.function.Consumer;

final class Dock extends LinearLayout {
    private static final String GRID_ASSET = "branding/grid.png";

    private final TileFactory tiles;
    private final Consumer<Intent> launcher;
    private final Runnable onOpenDrawer;
    private String shownSignature = "";

    Dock(Context context, TileFactory tiles, Consumer<Intent> launcher, Runnable onOpenDrawer) {
        super(context);
        this.tiles = tiles;
        this.launcher = launcher;
        this.onOpenDrawer = onOpenDrawer;
        setOrientation(LinearLayout.HORIZONTAL);
        setGravity(Gravity.CENTER);
        setClipChildren(false);
        setClipToPadding(false);
    }

    boolean show(List<LauncherApp> apps, String signature) {
        if (signature.equals(shownSignature)) {
            return false;
        }
        shownSignature = signature;
        removeAllViews();
        for (LauncherApp app : apps) {
            addView(tiles.appTile(app, () -> launcher.accept(app.launchIntent)));
        }
        addView(tiles.iconTile(AssetImages.bitmap(getContext(), GRID_ASSET), "Applis", "Toutes les applis", onOpenDrawer));
        return true;
    }

    void focusFirst() {
        if (getChildCount() > 0) {
            getChildAt(0).requestFocus();
        }
    }

    void setBlocked(boolean blocked) {
        setDescendantFocusability(blocked ? ViewGroup.FOCUS_BLOCK_DESCENDANTS : ViewGroup.FOCUS_AFTER_DESCENDANTS);
    }

    static FrameLayout.LayoutParams layoutParams(int screenHeight) {
        FrameLayout.LayoutParams params = new FrameLayout.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT, Gravity.BOTTOM);
        params.bottomMargin = screenHeight / 14;
        return params;
    }
}
