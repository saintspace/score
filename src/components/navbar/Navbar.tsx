import styles from './Navbar.module.css';
import { useDisclosure } from '@mantine/hooks';
import { Drawer } from '@mantine/core';
import { Link } from "react-router-dom";

const Navbar = (props: any) => {
    const [opened, { open, close }] = useDisclosure(false);
    const title = props.title == null ? "" : props.title;

    return (
        <>
            <nav className={styles.navbar} aria-label="Main navigation">
                <div className={styles.logo}>
                    <img className={styles.logoImage} src="/saintspace-crossball-logo-128.png" alt="logo" />
                </div>
                <div className={styles.navbarTitle}>{title}</div>
                <div className={styles.navigation}>
                    <button className={styles.menuButton} onClick={open}>
                        <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 50 50">
                            <g fill="none" stroke="#f7fbfd" stroke-width="3">
                                <line x1="10" y1="15" x2="40" y2="15"/>
                                <line x1="10" y1="25" x2="40" y2="25"/>
                                <line x1="10" y1="35" x2="40" y2="35"/>
                            </g>
                        </svg>
                    </button>
                </div>    
            </nav>
            <Drawer opened={opened} onClose={close} title="" position="right">
                <ul className={styles.navigationLinks}>
                    <li className={styles.navigationLink}>
                    <Link to={`/`} onClick={close} aria-label="Home">Home</Link>
                        <hr/>
                    </li>
                    <li className={styles.navigationLink}>
                        <Link to={`/saintspace/universe/account`} onClick={close} aria-label="Account">Account</Link>
                        <hr/>
                    </li>
                    <li className={styles.navigationLink}>
                        <Link to={`/saintspace/universe/auth/signout`} onClick={close} aria-label="Sign Out">Sign Out</Link>
                    </li>
                </ul>
            </Drawer>
        </>    
    );
}

export default Navbar;