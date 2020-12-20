<template>
    <v-navigation-drawer app>
        <v-list-item>
          <v-list-item-content>
            <v-list-item-title class="title">
              {{ $t('app.name')}} 
            </v-list-item-title>
            <v-list-item-subtitle>
              {{ $t('app.description')}}
            </v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
  
        <v-divider></v-divider>

        <v-list dense nav>
          <v-list-item
            v-for="item in items"
            :key="item.title"
            :to="item.target"
            link>
            <v-list-item-icon>
              <v-icon>{{ item.icon }}</v-icon>
            </v-list-item-icon>
  
            <v-list-item-content>
              <v-list-item-title>{{ item.title }}</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
    </v-navigation-drawer>
</template>

<script>
    export default {
        data() {
            return {
                items: [
                    { title: this.$i18n.t('navbar.title.Home'), icon: 'mdi-home', target: "/" },
                    { title: this.$i18n.t('navbar.title.Projects'), icon: 'mdi-image', target: "/project" },
                    { title: this.$i18n.t('navbar.title.Help'), icon: 'mdi-help-box', target: "/help" },
                    { title: this.$i18n.t('navbar.title.About'), icon: 'mdi-information', target: "/about" },
                ],
            }
        }
    }
</script>