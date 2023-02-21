import {Outlet} from "react-router-dom";
import Container from "react-bootstrap/Container";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import {memo} from "react";

function Navigation() {
    return (
        <Navbar bg="light" expand="lg">
            <Container>
                <Navbar.Brand href="/">Battleship</Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="me-auto">
                        <Nav.Item>
                            <Nav.Link href="/games">Games</Nav.Link>
                        </Nav.Item>
                        <Nav.Item>
                            <Nav.Link href="/scoreboard">Scoreboard</Nav.Link>
                        </Nav.Item>
                        <Nav.Item>
                            <Nav.Link href="/about">About</Nav.Link>
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
        <>
            <Header />
            <Navigation />
            <Outlet />
        </>
    )
}

export default memo(Layout)