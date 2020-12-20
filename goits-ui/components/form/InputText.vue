<template>
    <v-text-field
        v-model="modelValue"
        :rules="rules"
        :label="inputLabel"
        :required="required">
    </v-text-field>
</template>

<script>

import BasicField from '@/components/mixin/BasicField.js'
import FormField from '@/components/mixin/FormField.js'

export default {
    mixins: [BasicField, FormField],

    props: {
    }
}
</script>