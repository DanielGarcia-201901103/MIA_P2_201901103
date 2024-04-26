import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';
import House from '../House/House.jsx';
import Pantalla1 from '../Pantalla1/Pantalla1.jsx';
import Pantalla2 from '../Pantalla2/Pantalla2.jsx';
import Pantalla3 from '../Pantalla3/Pantalla3.jsx';

function Navb() {
  return (
    <>
      <Navbar expand="lg" bg="dark" data-bs-theme="dark">
        <Container>
          <Navbar.Brand href="/home"><House/>Home</Navbar.Brand>
          <Nav className="me-auto">
            <Nav.Link href="/Pantalla1"><Pantalla1/>Pantalla 1</Nav.Link>
            <Nav.Link href="/Pantalla2"><Pantalla2/>Pantalla 2</Nav.Link>
            <Nav.Link href="/Pantalla3"><Pantalla3/>Pantalla 3</Nav.Link>
          </Nav>
        </Container>
      </Navbar>
    </>
  );
}

export default Navb;