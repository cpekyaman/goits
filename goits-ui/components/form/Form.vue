<template>
    <v-form :ref="formRef"
        v-model="valid"
        lazy-validation>

        <slot></slot>

        <v-btn
            :disabled="!valid"
            color="success"
            class="mr-4"
            @click="onSubmit">
            {{ $t('form.button.submit') }}
        </v-btn>

        <v-btn
            v-if="showReset"
            color="error"
            class="mr-4"
            @click="reset">
            {{ $t('form.button.reset') }}
        </v-btn>
    </v-form>
</template>

<script>
export default {
    data() {
        return {
            valid: false
        }
    },

    props: {
        formRef: {type: String, required: true},
        showReset: {type: Boolean, default: false},
        onSubmit: {type: Function, required: true}
    },

    methods: {
        validate () {
            this.$refs[this.formRef].validate()
        },
        reset () {
            this.$refs[this.formRef].reset()
        }
    },
}
</script>