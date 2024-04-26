import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';
import { useNavigate } from 'react-router-dom';
import './Navb.css'

function Navb() {
  return (
    <>
      <Navbar expand="lg" bg="dark" data-bs-theme="dark"  >
        <Container>
          <Navbar.Brand href="/home">Home</Navbar.Brand>
          <Nav className="me-auto">
            <Nav.Link href="/Pantalla1">Pantalla 1</Nav.Link>
            <Nav.Link href="/Pantalla2">Pantalla 2</Nav.Link>
            <Nav.Link href="/Pantalla3">Pantalla 3</Nav.Link>
          </Nav>
        </Container>
      </Navbar>
    </>
  );
}

export default Navb;