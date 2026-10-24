package fr.sway.accueil;

import android.app.Activity;
import android.os.Bundle;

public class FireplaceActivity extends Activity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(new FireplaceView(this));
    }
}
