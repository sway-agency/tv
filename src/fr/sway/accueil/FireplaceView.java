package fr.sway.accueil;

import android.content.Context;
import android.content.res.AssetFileDescriptor;
import android.graphics.Color;
import android.media.AudioAttributes;
import android.media.MediaPlayer;
import android.view.Gravity;
import android.view.SurfaceHolder;
import android.view.SurfaceView;
import android.widget.FrameLayout;

import java.io.IOException;

final class FireplaceView extends FrameLayout implements SurfaceHolder.Callback {
    private static final String VIDEO_ASSET = "fireplace/fire.mp4";
    private static final String CRACKLE_ASSET = "fireplace/crackle.m4a";
    private static final float VIDEO_VOLUME = 0.6f;
    private static final float CRACKLE_VOLUME = 0.6f;

    private MediaPlayer videoPlayer;
    private MediaPlayer cracklePlayer;

    FireplaceView(Context context) {
        super(context);
        setBackgroundColor(Color.BLACK);
        setKeepScreenOn(true);
        SurfaceView surface = new SurfaceView(context);
        surface.getHolder().addCallback(this);
        addView(surface, new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT, Gravity.CENTER));
    }

    @Override
    public void surfaceCreated(SurfaceHolder holder) {
        videoPlayer = loopingPlayer(VIDEO_ASSET, VIDEO_VOLUME);
        if (videoPlayer != null) {
            videoPlayer.setDisplay(holder);
            videoPlayer.prepareAsync();
        }
        cracklePlayer = loopingPlayer(CRACKLE_ASSET, CRACKLE_VOLUME);
        if (cracklePlayer != null) {
            cracklePlayer.prepareAsync();
        }
    }

    @Override
    public void surfaceChanged(SurfaceHolder holder, int format, int width, int height) {
    }

    @Override
    public void surfaceDestroyed(SurfaceHolder holder) {
        release();
    }

    @Override
    protected void onDetachedFromWindow() {
        release();
        super.onDetachedFromWindow();
    }

    private MediaPlayer loopingPlayer(String asset, float volume) {
        MediaPlayer mediaPlayer = new MediaPlayer();
        mediaPlayer.setAudioAttributes(new AudioAttributes.Builder()
                .setUsage(AudioAttributes.USAGE_MEDIA)
                .setContentType(AudioAttributes.CONTENT_TYPE_MOVIE)
                .build());
        mediaPlayer.setLooping(true);
        mediaPlayer.setVolume(volume, volume);
        mediaPlayer.setOnPreparedListener(MediaPlayer::start);
        try (AssetFileDescriptor descriptor = getContext().getAssets().openFd(asset)) {
            mediaPlayer.setDataSource(descriptor.getFileDescriptor(), descriptor.getStartOffset(), descriptor.getLength());
            return mediaPlayer;
        } catch (IOException e) {
            mediaPlayer.release();
            return null;
        }
    }

    private void release() {
        videoPlayer = release(videoPlayer);
        cracklePlayer = release(cracklePlayer);
    }

    private static MediaPlayer release(MediaPlayer player) {
        if (player != null) {
            player.release();
        }
        return null;
    }
}
