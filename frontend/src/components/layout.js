import {Link, Outlet} from "react-router-dom";
import Container from "react-bootstrap/Container";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import {NavLink} from "react-bootstrap";

function Navigation() {
    return (
        <Navbar bg="light" expand="lg">
            <Container>
                <Navbar.Brand href="/">Battleship</Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="me-auto">
                        <Nav.Item>
                            <NavLink as={Link} to="/games">Games</NavLink>
                        </Nav.Item>
                        <Nav.Item>
                            <NavLink as={Link} to="/scoreboard">Scoreboard</NavLink>
                        </Nav.Item>
                        <Nav.Item>
                            <NavLink as={Link} to="/about">About</NavLink>
                        </Nav.Item>
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    )
}

function Header() {
    return (
        <Container>
            <h1 className="text-center">Let's play 'Battleship'</h1>
        </Container>
    );
}

const Layout = () => {
    return (
        <Container>
            <Header />
            <Navigation />
            <Outlet />
        </Container>
    )
}

export default Layout