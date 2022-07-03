import { AppBar, Menu, MenuItem, Toolbar, Typography } from "@mui/material";
import { useContext, useState } from "react";
import { useQueryClient } from "react-query";
import { useNavigate } from "react-router-dom";
import { authContext, logout } from "../../api/auth";
import { useUsername } from "../../api/profile";

export default function Navbar() {
  const [usernameEl, setUsernameEl] = useState<null | HTMLElement>(null);
  const username = useUsername();
  const navigate = useNavigate();
  const client = useQueryClient();
  const { setAuthenticated } = useContext(authContext);
  async function doLogout() {
    logout();
    setAuthenticated(false);
    await client.removeQueries("assessments")
  }
  return (
    <AppBar position="relative">
      <Toolbar>
        <Typography sx={{ flexGrow: 1, cursor: 'pointer' }} onClick={() => navigate("/")}>Digital Dexterity Measurement</Typography>
        <Typography
          sx={{ cursor: 'pointer' }}
          onClick={(e) => setUsernameEl(e.currentTarget)}
        >{username}</Typography>
        <Menu
          anchorEl={usernameEl}
          open={Boolean(usernameEl)}
          onClose={() => setUsernameEl(null)}
        >
          <MenuItem sx={{ width: 120 }} onClick={() => doLogout()}>Logout</MenuItem>
        </Menu>
      </Toolbar>
    </AppBar>
  )
}