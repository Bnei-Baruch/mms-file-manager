import React, {PropTypes} from 'react'

export const VisibleItemList = (props) => (
    <ul>
        {props.items.map(item => (
            <li onClick={e => props.onSelected(item.id)} key={item.id}>{item.name}</li>
        ))}
    </ul>
);

VisibleItemList.propTypes = {
    items: PropTypes.arrayOf(PropTypes.shape({
        id: PropTypes.number.isRequired,
        name: PropTypes.string.isRequired
    })),
    onSelected: PropTypes.func.isRequired
};

export default VisibleItemList;

