package fr.sway.accueil;

import android.content.Context;
import android.graphics.Color;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.GridLayout;
import android.widget.LinearLayout;
import android.widget.ScrollView;
import android.widget.TextView;

import java.util.List;

final class AppDrawer extends ScrollView {
    private static final int COLUMNS = 5;
    private static final long FADE_MILLIS = 200;

    private final GridLayout grid;
    private final int screenHeight;

    AppDrawer(Context context, int screenHeight) {
        super(context);
        this.screenHeight = screenHeight;
        setBackgroundColor(Color.argb(225, 10, 10, 14));
        setElevation(screenHeight);
        setVisibility(View.GONE);
        setClipChildren(false);

        LinearLayout content = new LinearLayout(context);
        content.setOrientation(LinearLayout.VERTICAL);
        content.setGravity(Gravity.CENTER_HORIZONTAL);
        content.setClipChildren(false);
        content.setClipToPadding(false);
        content.setPadding(0, screenHeight / 12, 0, screenHeight / 10);
        content.addView(title(context), titleParams());

        grid = new GridLayout(context);
        grid.setColumnCount(COLUMNS);
        grid.setClipChildren(false);
        grid.setClipToPadding(false);
        content.addView(grid, new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT));
        addView(content, new ScrollView.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT));
    }

    boolean isOpen() {
        return getVisibility() == View.VISIBLE;
    }

    void open(List<View> tiles) {
        grid.removeAllViews();
        for (View tile : tiles) {
            grid.addView(tile);
        }
        setVisibility(View.VISIBLE);
        setAlpha(0f);
        animate().alpha(1f).setDuration(FADE_MILLIS).start();
        scrollTo(0, 0);
        Motion.riseIn(grid, screenHeight / 20f);
        if (grid.getChildCount() > 0) {
            grid.getChildAt(0).requestFocus();
        }
    }

    void close() {
        setVisibility(View.GONE);
    }

    private TextView title(Context context) {
        TextView title = new TextView(context);
        title.setText("Applications");
        title.setTextColor(Color.WHITE);
        title.setTypeface(Fonts.satoshi(context.getAssets(), 300));
        title.setTextSize(TypedValue.COMPLEX_UNIT_PX, screenHeight / 22f);
        return title;
    }

    private LinearLayout.LayoutParams titleParams() {
        LinearLayout.LayoutParams params = new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT);
        params.bottomMargin = screenHeight / 30;
        return params;
    }
}
