import React from 'react';
import VisibleItemList from './VisibleItemList'
import FilterList from './FilterList'

export const Inbox = (props) => (
    <div>
        <h1>Inbox</h1>
        <VisibleItemList items={props.items} onSelected={props.onItemSelected}/>
        <FilterList active={props.visibilityFilter} onSelected={props.onFilterSelected}/>
    </div>
);

export default Inbox;
