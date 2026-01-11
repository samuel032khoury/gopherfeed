import { useNavigate } from "react-router-dom";
import { BASE_URL } from "./constant";


export const ActivationPage = () => {
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token') || '';
    const redirect = useNavigate();
    const handleConfirm = async (token: string) => {

        const response = await fetch(`${BASE_URL}/auth/activate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({token}),
        });
        if (response.ok) {
            redirect('/');
        } else {
            alert('There was an error confirming your registration.');
        }
    }
    return (
        <div>
            <h1>Registration Confirmed</h1>
            <button onClick={() => handleConfirm(token)}>Click to confirm</button>
        </div>
    )};