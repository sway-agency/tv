package fr.sway.accueil;

import android.view.View;
import android.view.ViewGroup;

final class Motion {
    static final SpringInterpolator SPRING = SpringInterpolator.gentle();
    static final long SPRING_MILLIS = 520;

    private static final long ENTRANCE_MILLIS = 700;
    private static final long ENTRANCE_STAGGER_MILLIS = 55;

    private Motion() {
    }

    static void springScale(View view, float scale) {
        view.animate().scaleX(scale).scaleY(scale).setInterpolator(SPRING).setDuration(SPRING_MILLIS).start();
    }

    static void riseIn(ViewGroup group, float distance) {
        for (int i = 0; i < group.getChildCount(); i++) {
            View child = group.getChildAt(i);
            child.setAlpha(0f);
            child.setTranslationY(distance);
            child.animate().alpha(1f).translationY(0f)
                    .setStartDelay(i * ENTRANCE_STAGGER_MILLIS)
                    .setInterpolator(SPRING)
                    .setDuration(ENTRANCE_MILLIS)
                    .start();
        }
    }
}
