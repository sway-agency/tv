package fr.sway.accueil;

import android.content.Context;
import android.graphics.Color;
import android.graphics.Typeface;
import android.graphics.drawable.GradientDrawable;
import android.util.TypedValue;
import android.view.Gravity;
import android.view.View;
import android.widget.ImageView;
import android.widget.LinearLayout;
import android.widget.TextClock;

final class TopBar {
    private static final String SWAY_LOGO_ASSET = "branding/sway_s.png";
    private static final String GEAR_ASSET = "branding/gear.png";

    private final Context context;
    private final int screenWidth;
    private final int screenHeight;

    TopBar(Context context, int screenWidth, int screenHeight) {
        this.context = context;
        this.screenWidth = screenWidth;
        this.screenHeight = screenHeight;
    }

    View clock() {
        LinearLayout clock = row(LinearLayout.VERTICAL, Gravity.START);
        clock.addView(clockText("HH:mm", screenHeight / 11f, Fonts.satoshi(context.getAssets(), 300)));
        clock.addView(clockText("EEEE d MMMM", screenHeight / 40f, Fonts.satoshi(context.getAssets(), 400)));
        return clock;
    }

    View corner(Runnable onSettings) {
        LinearLayout corner = row(LinearLayout.HORIZONTAL, Gravity.CENTER_VERTICAL);
        int gearSize = screenHeight / 14;
        LinearLayout.LayoutParams gearParams = new LinearLayout.LayoutParams(gearSize, gearSize);
        gearParams.rightMargin = screenHeight / 40;
        corner.addView(settingsButton(gearSize, onSettings), gearParams);

        int logoSize = screenHeight / 24;
        ImageView logo = AssetImages.view(context, SWAY_LOGO_ASSET);
        logo.setAlpha(0.85f);
        corner.addView(logo, new LinearLayout.LayoutParams(logoSize, logoSize));
        return corner;
    }

    private LinearLayout row(int orientation, int gravity) {
        LinearLayout row = new LinearLayout(context);
        row.setOrientation(orientation);
        row.setGravity(gravity);
        row.setClipChildren(false);
        row.setClipToPadding(false);
        int margin = screenWidth / 32;
        row.setPadding(margin, margin * 3 / 4, margin, 0);
        return row;
    }

    private TextClock clockText(String pattern, float sizePx, Typeface typeface) {
        TextClock text = new TextClock(context);
        text.setFormat12Hour(pattern);
        text.setFormat24Hour(pattern);
        text.setTextColor(Color.WHITE);
        text.setTextSize(TypedValue.COMPLEX_UNIT_PX, sizePx);
        text.setTypeface(typeface);
        text.setShadowLayer(12f, 0f, 2f, Color.argb(160, 0, 0, 0));
        return text;
    }

    private View settingsButton(int size, Runnable onSettings) {
        ImageView gear = AssetImages.view(context, GEAR_ASSET);
        int inset = size / 4;
        gear.setPadding(inset, inset, inset, inset);
        GradientDrawable highlight = new GradientDrawable();
        highlight.setShape(GradientDrawable.OVAL);
        highlight.setColor(Color.argb(70, 255, 255, 255));
        highlight.setAlpha(0);
        gear.setBackground(highlight);
        gear.setContentDescription("Réglages");
        gear.setFocusable(true);
        gear.setOnFocusChangeListener((view, focused) -> {
            highlight.setAlpha(focused ? 255 : 0);
            Motion.springScale(view, focused ? 1.15f : 1f);
        });
        gear.setOnClickListener(view -> onSettings.run());
        return gear;
    }
}
