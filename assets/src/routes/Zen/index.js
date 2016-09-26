import {injectReducer} from '../../store/reducers'

export default (store) => ({
    path: '/zen',
    getComponent(nextState, cb) {
        require.ensure([], (require) => {
            const container = require('./containers/ZenContainer').default,
                reducer = require('./modules/zen').default;
            injectReducer(store, {key: 'zen', reducer});
            cb(null, container)
        }, 'zen')
    }
})
