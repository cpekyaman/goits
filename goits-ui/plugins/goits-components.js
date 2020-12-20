import Vue from 'vue'

// global layout components
import TopBar from '@/components/layout/TheTopBar.vue'
import NavBar from '@/components/layout/TheNavBar.vue'
import SearchBar from '@/components/layout/TheSearchBar.vue'

Vue.component('g-top-bar', TopBar)
Vue.component('g-nav-bar', NavBar)
Vue.component('g-search-bar', SearchBar)

// page components
import EntityList from '@/components/page/EntityList.vue'
import EntityDetail from '@/components/page/EntityDetail.vue'
Vue.component('g-entity-list', EntityList)
Vue.component('g-entity-detail', EntityDetail)

// form components
import Form from '@/components/form/Form.vue'
import InputPassword from '@/components/form/InputPassword.vue'
import InputText from '@/components/form/InputText.vue'

Vue.component('g-form', Form)
Vue.component('g-input-password', InputPassword)
Vue.component('g-input-text', InputText)

// display components
import DisplayText from '@/components/display/DisplayText.vue'
Vue.component('g-display-text', DisplayText)