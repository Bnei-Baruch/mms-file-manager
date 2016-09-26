import React from 'react'

export const FileDetails = (props) => (
    <div>
        <ul>
            <li>Details</li>
            <li>Validations</li>
            <li>History</li>
        </ul>
        <div>{props.file.name}: {props.file.id}</div>
    </div>
);

export default FileDetails;