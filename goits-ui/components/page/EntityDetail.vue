<template>
    <v-row class="detail-page entity-detail">
        <v-col cols="12">
            <slot></slot>
        </v-col>
    </v-row>
</template>

<script>

import { mapMutations } from 'vuex'
import BasePageComponent from '@/components/mixin/BasePageComponent.js'
import RouteAware from '@/components/mixin/RouteAware.js'

export default {
    mixins: [BasePageComponent,RouteAware],

    props : {
        itemId: {type: String, required: true}
    },

    async fetch() {
        const { data } = await this.$axios.get(this.rootPath + '/' + this.itemId)
        this.$store.commit(this.rootPath + '/item', data.data)
    }
}
</script>