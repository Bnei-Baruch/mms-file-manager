import React from 'react'
import {IndexLink, Link} from 'react-router'
import classes from './Header.scss'

export const Header = () => (
    <div>
        <h1>BB Multimedia Management System</h1>
        <IndexLink to='/' activeClassName={classes.activeRoute}>
            Home
        </IndexLink>
        {' · '}
        <Link to='/counter' activeClassName={classes.activeRoute}>
            Counter
        </Link>
        {' · '}
        <Link to='/zen' activeClassName={classes.activeRoute}>
            Zen
        </Link>
        {' · '}
        <Link to='/inbox' activeClassName={classes.activeRoute}>
            Inbox
        </Link>
    </div>
);

export default Header
