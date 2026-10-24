package fr.sway.accueil;

import android.service.dreams.DreamService;

public class FireplaceDream extends DreamService {
    @Override
    public void onAttachedToWindow() {
        super.onAttachedToWindow();
        setInteractive(false);
        setFullscreen(true);
        setScreenBright(true);
        setContentView(new FireplaceView(this));
    }
}
