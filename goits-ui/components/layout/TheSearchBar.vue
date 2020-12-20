<template>
    <v-app-bar app>
        <v-app-bar-nav-icon></v-app-bar-nav-icon>
  
        <v-toolbar-title>{{ title }}</v-toolbar-title>

        <v-spacer></v-spacer>
  
        <v-btn icon>
          <v-icon>mdi-magnify</v-icon>
        </v-btn>
    </v-app-bar>
</template>

<script>
export default {
    props: {
        title: {type: String, default: "Hello"}
    }
}
</script>