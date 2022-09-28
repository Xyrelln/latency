import Vuex from 'vuex'
import user, { User } from './modules/user'

export interface StoreState {
  user: User
}


export const store = new Vuex.Store<StoreState>({
  modules: {
    user
  },
})