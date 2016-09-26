import React, {PropTypes} from 'react'
import classNames from 'classnames'
import classes from './FilterList.scss'

export const FilterList = (props) => (
    <ul>
        <li className={classNames({[classes.active]: props.active === 'SHOW_ALL'})}
            onClick={e => props.onSelected('SHOW_ALL')}>All</li>
        <li className={classNames({[classes.active]: props.active === 'SHOW_COMPLETE'})}
            onClick={e => props.onSelected('SHOW_COMPLETE')}>Complete</li>
        <li className={classNames({[classes.active]: props.active === 'SHOW_IN_PROGRESS'})}
            onClick={e => props.onSelected('SHOW_IN_PROGRESS')}>In Progress</li>
    </ul>
);

FilterList.propTypes = {
    onSelected: PropTypes.func.isRequired,
    active: PropTypes.string
};

export default FilterList;

