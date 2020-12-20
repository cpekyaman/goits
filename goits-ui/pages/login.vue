<template>
    <v-row class="form-page login-form">
        <v-col cols="8">
            <g-form formRef="loginForm" :onSubmit="doLogin">
                <g-input-text :dataModel="user" field="username" :rules="usernameRules" />
                <g-input-password :dataModel="user" field="password" :rules="passwordRules" />
            </g-form>
        </v-col>
    </v-row>
</template>

<script>
export default {
    data() {
        return {
            valid: false,
            user: {
                username: '',
                password: ''
            },
            showPassword: false,
            usernameRules: [
                v => !!v || 'Username is required',
                v => (v && v.length <= 10) || 'Username must be less than 10 characters',
            ],
            passwordRules: [
                v => !!v || 'Password is required',
                v => (v && v.length >= 10 && v.length <= 25) || 'Password must be less than 10 characters',
            ],
        }
    },

    methods: {
        doLogin() {
            this.$axios.$post('login', {
                username: 'cenk',
                password: 'password'
            })    
        }
    }
}
</script>