package fr.sway.accueil;

import android.animation.TimeInterpolator;

final class SpringInterpolator implements TimeInterpolator {
    private final float damping;
    private final float omega;
    private final float dampedOmega;

    SpringInterpolator(float damping, float omega) {
        this.damping = damping;
        this.omega = omega;
        this.dampedOmega = omega * (float) Math.sqrt(1 - damping * damping);
    }

    static SpringInterpolator gentle() {
        return new SpringInterpolator(0.5f, 12f);
    }

    @Override
    public float getInterpolation(float t) {
        float decay = (float) Math.exp(-damping * omega * t);
        float phase = dampedOmega * t;
        return 1f - decay * (float) (Math.cos(phase) + damping * omega / dampedOmega * Math.sin(phase));
    }
}
