import React , {Component, PropTypes} from 'react'
import {connect} from 'react-redux'
import {fetchInboxItems, filterChanged} from '../modules/inbox'
import Inbox from '../components/Inbox'
import {getVisibleItems} from '../selectors'

class InboxContainer extends Component {

    static propTypes = {
        fetchInboxItems: PropTypes.func.isRequired,
        filterChanged: PropTypes.func.isRequired
    };

    componentDidMount() {
        this.props.fetchInboxItems();
    }

    handleFilterSelected(newFilter) {
        this.props.fetchInboxItems();

        // TODO: should be on Promise.then()
        this.props.filterChanged(newFilter);
    }

    handleItemSelected(itemId) {
        console.log('item selected', itemId);
        this.props.history.push(`/files/${itemId}`);
    }

    render() {
        return (<Inbox {...this.props}
                       onFilterSelected={(x) => this.handleFilterSelected(x)}
                       onItemSelected={(x) => this.handleItemSelected(x)}/>)
    }
}

const mapDispatchToProps = {
    fetchInboxItems,
    filterChanged
};

const mapStateToProps = (state) => {
    return {
        items: getVisibleItems(state.inbox),
        visibilityFilter: state.inbox.visibilityFilter
    }
};

export default connect(mapStateToProps, mapDispatchToProps)(InboxContainer)