import {injectReducer} from '../../../../store/reducers'

export default (store) => ({
    path: '/files/:id',
    getComponent(nextState, cb) {
        require.ensure([], (require) => {
            const container = require('./containers/FileDetailsContainer').default,
                reducer = require('./modules/fileDetails').default;
            injectReducer(store, {key: 'inbox.details', reducer});
            cb(null, container);
        }, 'fileDetails')
    }
})