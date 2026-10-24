package fr.sway.accueil;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.Color;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.view.ViewGroup;
import android.widget.FrameLayout;
import android.widget.GridLayout;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.ScrollView;
import android.widget.TextView;

import java.util.List;

final class AppDrawer extends FrameLayout {
    private static final int COLUMNS = 5;
    private static final long FADE_MILLIS = 200;
    private static final int TINT = Color.argb(70, 22, 22, 26);

    private final ImageView backdrop;
    private final ScrollView scroll;
    private final GridLayout grid;
    private final int screenHeight;

    AppDrawer(Context context, int screenHeight) {
        super(context);
        this.screenHeight = screenHeight;
        setElevation(screenHeight);
        setVisibility(View.GONE);
        setBackgroundColor(Color.rgb(18, 18, 21));

        backdrop = new ImageView(context);
        backdrop.setScaleType(ImageView.ScaleType.CENTER_CROP);
        addView(backdrop, fill());
        View tint = new View(context);
        tint.setBackgroundColor(TINT);
        addView(tint, fill());
        View grain = new View(context);
        grain.setBackground(Grain.drawable(getResources()));
        addView(grain, fill());

        grid = new GridLayout(context);
        grid.setColumnCount(COLUMNS);
        grid.setClipChildren(false);
        grid.setClipToPadding(false);
        scroll = new ScrollView(context);
        scroll.setFillViewport(true);
        scroll.setClipChildren(false);
        scroll.addView(centeredContent(context), new ScrollView.LayoutParams(
                ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.WRAP_CONTENT));
        addView(scroll, fill());
    }

    boolean isOpen() {
        return getVisibility() == View.VISIBLE;
    }

    void open(List<View> tiles, Bitmap backdropImage) {
        backdrop.setImageBitmap(backdropImage);
        grid.removeAllViews();
        for (View tile : tiles) {
            grid.addView(tile);
        }
        setVisibility(View.VISIBLE);
        setAlpha(0f);
        animate().alpha(1f).setDuration(FADE_MILLIS).start();
        scroll.scrollTo(0, 0);
        Motion.riseIn(grid, screenHeight / 20f);
        if (grid.getChildCount() > 0) {
            grid.getChildAt(0).requestFocus();
        }
    }

    void close() {
        setVisibility(View.GONE);
        backdrop.setImageDrawable(null);
    }

    private LinearLayout centeredContent(Context context) {
        LinearLayout content = new LinearLayout(context);
        content.setOrientation(LinearLayout.VERTICAL);
        content.setGravity(Gravity.CENTER);
        content.setClipChildren(false);
        content.setClipToPadding(false);
        content.setPadding(0, screenHeight / 12, 0, screenHeight / 12);
        LinearLayout.LayoutParams titleParams = new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT);
        titleParams.bottomMargin = screenHeight / 30;
        content.addView(title(context), titleParams);
        content.addView(grid, new LinearLayout.LayoutParams(
                ViewGroup.LayoutParams.WRAP_CONTENT, ViewGroup.LayoutParams.WRAP_CONTENT));
        return content;
    }

    private TextView title(Context context) {
        TextView title = new TextView(context);
        title.setText("Applications");
        title.setTextColor(Color.WHITE);
        title.setTypeface(Fonts.satoshi(context.getAssets(), 700));
        title.setTextSize(TypedValue.COMPLEX_UNIT_PX, screenHeight / 22f);
        return title;
    }

    private static FrameLayout.LayoutParams fill() {
        return new FrameLayout.LayoutParams(ViewGroup.LayoutParams.MATCH_PARENT, ViewGroup.LayoutParams.MATCH_PARENT);
    }
}
