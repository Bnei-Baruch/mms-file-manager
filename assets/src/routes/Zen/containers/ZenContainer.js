import {connect} from 'react-redux'
import {fetchZen, saveCurrentZen} from '../modules/zen'
import Zen from '../components/Zen'

const mapDispatchToProps = {
    fetchZen,
    saveCurrentZen
};

const mapStateToProps = (state) => ({
    zen: state.zen.zens.find(x => x.id === state.zen.current),
    saved: state.zen.zens.filter(x => state.zen.saved.indexOf(x.id) !== -1)
});

export default connect(mapStateToProps, mapDispatchToProps)(Zen)