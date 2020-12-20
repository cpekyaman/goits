<template>
    <v-system-bar app
        height="30" 
        dark>

        <v-spacer></v-spacer>

        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-btn icon 
              v-bind="attrs"
              v-on="on"
              to="/login">
              <v-icon>mdi-login</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('statusbar.button.login') }}</span>
        </v-tooltip>

        <v-tooltip bottom>
          <template v-slot:activator="{ on, attrs }">
            <v-btn icon 
              v-bind="attrs"
              v-on="on">
              <v-icon>mdi-account-plus</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('statusbar.button.signup') }}</span>
        </v-tooltip>
    </v-system-bar>
</template>