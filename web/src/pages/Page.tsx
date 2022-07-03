import { Container } from "@mui/material";
import { ReactNode } from "react";
import Navbar from "../components/shared/Navbar";

export default function Page(props: { to: ReactNode }) {
  return (
    <>
      <Navbar />
      <Container>
        {props.to}
      </Container>
    </>
  )
}