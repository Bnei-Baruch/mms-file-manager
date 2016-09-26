import {injectReducer} from '../../store/reducers'
import EmptyLayout from '../../layouts/EmptyLayout'
import FileDetailsRoute from './routes/Details'
import InboxContainer from './containers/InboxContainer'
import reducer from './modules/inbox'

export default (store) => {
    injectReducer(store, {key: 'inbox', reducer});
    return {
        path: '/inbox',
        component: EmptyLayout,
        childRoutes: [
            FileDetailsRoute(store)
        ],
        indexRoute: {component: InboxContainer}
    }
}