export const state = () => ({
    user: {
        userName: '',
        fullName: '',
        email: ''
    }
})
  
export const mutations = {
    login(state, { user }) {
        state.user = user
    }
}
  