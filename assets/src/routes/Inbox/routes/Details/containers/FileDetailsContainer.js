import React, {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import FileDetails from '../components/FileDetails'

class FileDetailsContainer extends Component {

    componentDidMount() {
        console.log('FileDetails did mount')
    }

    render() {
        return (<FileDetails {...this.props}/>)
    }
}

const mapDispatchToProps = {};

const mapStateToProps = (state, ownProps) => {
    return {
        file: state.inbox.items.find(x => x.id == ownProps.routeParams.id )
    }
};

export default connect(mapStateToProps, mapDispatchToProps)(FileDetailsContainer);