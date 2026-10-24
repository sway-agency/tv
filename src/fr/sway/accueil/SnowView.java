package fr.sway.accueil;

import android.content.Context;
import android.graphics.Canvas;
import android.graphics.Color;
import android.graphics.Paint;
import android.view.View;

import java.util.Random;

final class SnowView extends View {
    private static final int FLAKE_COUNT = 70;

    private final Random random = new Random();
    private final Paint paint = new Paint(Paint.ANTI_ALIAS_FLAG);
    private final float[] x = new float[FLAKE_COUNT];
    private final float[] y = new float[FLAKE_COUNT];
    private final float[] radius = new float[FLAKE_COUNT];
    private final float[] fallSpeed = new float[FLAKE_COUNT];
    private final float[] swayPhase = new float[FLAKE_COUNT];
    private final int[] alpha = new int[FLAKE_COUNT];
    private long lastFrameNanos;
    private float elapsedSeconds;
    private boolean running;

    SnowView(Context context) {
        super(context);
        paint.setColor(Color.WHITE);
    }

    void start() {
        running = true;
        lastFrameNanos = 0;
        postInvalidateOnAnimation();
    }

    void stop() {
        running = false;
    }

    @Override
    protected void onSizeChanged(int width, int height, int oldWidth, int oldHeight) {
        for (int i = 0; i < FLAKE_COUNT; i++) {
            respawn(i, width, height, true);
        }
    }

    @Override
    protected void onDraw(Canvas canvas) {
        int width = getWidth();
        int height = getHeight();
        long now = System.nanoTime();
        float delta = lastFrameNanos == 0 ? 0f : Math.min((now - lastFrameNanos) / 1e9f, 0.05f);
        lastFrameNanos = now;
        elapsedSeconds += delta;

        float gust = (float) Math.sin(elapsedSeconds * 0.13f) * height * 0.012f;
        for (int i = 0; i < FLAKE_COUNT; i++) {
            y[i] += fallSpeed[i] * height * delta;
            x[i] += (gust + (float) Math.sin(elapsedSeconds * 0.9f + swayPhase[i]) * radius[i] * 6f) * delta;
            if (y[i] - radius[i] > height || x[i] < -radius[i] * 4 || x[i] > width + radius[i] * 4) {
                respawn(i, width, height, false);
            }
            paint.setAlpha(alpha[i]);
            canvas.drawCircle(x[i], y[i], radius[i], paint);
        }
        if (running) {
            postInvalidateOnAnimation();
        }
    }

    private void respawn(int i, int width, int height, boolean anywhere) {
        float depth = random.nextFloat();
        radius[i] = height * (0.0012f + depth * 0.0028f);
        fallSpeed[i] = 0.02f + depth * 0.045f;
        alpha[i] = (int) (60 + depth * 140);
        swayPhase[i] = random.nextFloat() * 6.28f;
        x[i] = random.nextFloat() * width;
        y[i] = anywhere ? random.nextFloat() * height : -radius[i] * 2;
    }
}
