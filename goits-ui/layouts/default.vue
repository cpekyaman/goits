<template>
    <v-app>
        <g-top-bar />
        <g-nav-bar />
        <g-search-bar />

        <!-- Sizes your content based upon application components -->
        <v-main>
            <!-- Provides the application the proper gutter -->
            <v-container fluid class="main-content container">
                <nuxt />
            </v-container>
        </v-main>

        <v-footer app>
            <!-- -->
        </v-footer>
    </v-app>
</template>

<script>
export default {
    
}
</script>