import React from 'react'

export const EmptyLayout = ({children}) => (
    <div>{children}</div>
);

EmptyLayout.propTypes = {
    children: React.PropTypes.element.isRequired
};

export default EmptyLayout
