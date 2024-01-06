import { Client } from "@/lib/pb/common";
import { ClientItem } from "./client_item";

export interface ClientListProps {
	Clients: Client[]
}
export const ClientList: React.FC<ClientListProps> = ({ Clients }) => {
	return (
		<>
			<ClientItem Client={Clients[0]} />
		</>
	)
};