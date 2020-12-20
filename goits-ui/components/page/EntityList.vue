<template>
    <v-row class="list-page entity-list">
        <v-col cols="12">
            <v-data-table app
                :headers="headerList"
                :items="dataList"
                :items-per-page="5"
                item-key="itemKey">

                <template v-slot:item.actions="{ item }">
                    <v-icon small class="mr-2" @click="showItem(item)">
                        mdi-eye
                    </v-icon>
                    <v-icon small class="mr-2" @click="editItem(item)">
                        mdi-pencil
                    </v-icon>
                    <v-icon small @click="removeItem(item)">
                        mdi-delete
                    </v-icon>
                </template>
            </v-data-table>
        </v-col>
    </v-row>
</template>

<script>
import BasePageComponent from '@/components/mixin/BasePageComponent.js'
import RouteAware from '@/components/mixin/RouteAware.js'

export default {
    mixins: [BasePageComponent,RouteAware],

    data() {
        return {
            dataList: []
        }
    },

    props : {
        columns : {type : Array, required : true},
        itemKey: {type : String, required : true},
        editable: {type: Boolean, default: true},
        removeable: {type: Boolean, default: false}
    },

    computed: {
        headers() {
            let headers = []
            columns.forEach(col => {
                headers.push({
                    "value" : col,
                    "text" : this.getLabel('table.header', col)
                })
            })

            if(this.editable || this.removeable) {
                headers.push({ text: this.getLabel('table.header', 'actions'), value: 'actions', sortable: false })
            } 
            return headers
        }
    },

    methods: {
        showItem(item) {
            console.log('show item ' + item[this.itemKey])
            this.$router.push({ path: '/' + this.rootPath + '/' + item[this.itemKey] + '/show' })
        },
        editItem(item) {

        },
        removeItem(item) {

        }
    },

    async fetch() {
        const { data } = await this.$axios.get(this.rootPath)
        console.log(data)
        this.dataList = data.data
    }
}
</script>