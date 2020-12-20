<template>
    <v-text-field
        :value="dataValue"
        :label="inputLabel"
        :readOnly="true"
        outlined>
    </v-text-field>
</template>

<script>
import BasicField from '@/components/mixin/BasicField.js'

export default {
    mixins: [BasicField],

    props: {
    },

    computed: {
        dataValue() {
            return this.dataModel[this.field]
        }
    }
}
</script>