package fr.sway.accueil;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;

public class DockReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        new DockConfig(context).applyIntent(intent);
    }
}
