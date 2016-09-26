import {createSelector} from 'reselect'

const getVisibiltyFilter = (state) => state.visibilityFilter;
const getItems = (state) => state.items;

export const getVisibleItems = createSelector(
    [getVisibiltyFilter, getItems],
    (filter, items) => {
        switch (filter) {
            case 'SHOW_ALL':
                return items;
            case 'SHOW_COMPLETE':
                return items.filter(item => item.status === 'complete');
            case 'SHOW_IN_PROGRESS':
                return items.filter(item => item.status === 'in_progress');
            
        }
    }
);
