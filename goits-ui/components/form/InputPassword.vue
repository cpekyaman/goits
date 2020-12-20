<template>
    <v-text-field
        v-model="modelValue"
        :rules="rules"
        @click:append="showPassword = !showPassword"
        :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
        :type="showPassword ? 'text' : 'password'"
        :label="inputLabel"
        :required="required"
    >
    </v-text-field>
</template>

<script>

import BasicField from '@/components/mixin/BasicField.js'
import FormField from '@/components/mixin/FormField.js'

export default {
    mixins: [BasicField, FormField],

    data() {
        return {
            showPassword: false
        }
    }, 

    props: {
        
    }
}
</script>