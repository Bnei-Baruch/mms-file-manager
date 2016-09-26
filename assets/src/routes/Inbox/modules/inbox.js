// ------------------------------------
// Constants
// ------------------------------------
export const REQUEST_INBOX_ITEMS = 'REQUEST_INBOX_ITEMS';
export const RECEIVE_INBOX_ITEMS = 'RECEIVE_INBOX_ITEMS';
export const FILTER_CHANGED = 'FILTER_CHANGED';

// ------------------------------------
// Actions
// ------------------------------------
export function requestInboxItems () {
    return {
        type: REQUEST_INBOX_ITEMS
    }
}

export function receiveInboxItems(items) {
    return {
        type: RECEIVE_INBOX_ITEMS,
        payload: {
            items: items
        }
    }
}

export const fetchInboxItems = () => {
    return (dispatch) => {
        dispatch(requestInboxItems());

        // Here should come async server fetch stuff
        let files = [];
        for (let i=0; i<150; i++) {
            files.push({
                id: i,
                name: `File ${i}`,
                status: Math.random() > 0.5 ? 'complete' : 'in_progress'
            })
        }

        dispatch(receiveInboxItems(files));
    }
};

export function filterChanged(newFilter) {
    return {
        type: FILTER_CHANGED,
        payload: {newFilter}
    }
}

export const actions = [
    requestInboxItems,
    receiveInboxItems,
    fetchInboxItems,
    filterChanged
];

// ------------------------------------
// Action Handlers
// ------------------------------------
const ACTION_HANDLERS = {
    [REQUEST_INBOX_ITEMS]: (state) => ({ ...state, fetching: true }),
    [RECEIVE_INBOX_ITEMS]: (state, action) => ({
        ...state,
        fetching: false,
        items: action.payload.items
    }),
    [FILTER_CHANGED]: (state, action) => ({
        ...state,
        visibilityFilter: action.payload.newFilter
    })
};

// ------------------------------------
// Reducer
// ------------------------------------
const initialState = {items:[], visibilityFilter: 'SHOW_ALL'};
export default function inboxReducer(state = initialState, action) {
    const handler = ACTION_HANDLERS[action.type];
    return handler ? handler(state, action) : state;
}