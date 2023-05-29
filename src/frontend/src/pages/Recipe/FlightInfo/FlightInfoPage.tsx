import { VStack, Text, Box } from "@chakra-ui/react"
import GetFlight from "postAPI/flights/Get"
import React from "react"
import { NavigateFunction, Params } from "react-router-dom"
import { Flight as FlightT} from "types/Flight"

type State = {
    flight?: FlightT
}

type RecipeInfoParams = {
    match: Readonly<Params<string>>
    navigate: NavigateFunction
}

class FlightInfoPage extends React.Component<RecipeInfoParams, State> {
    flightNumber: string;

    constructor(props) {
        super(props);
        this.flightNumber = this.props.match.flightNumber || "?";
        this.state = {}
    };

    componentDidMount(): void {
        GetFlight(this.flightNumber).then(data => {
            console.log(data)
            if (data.status === 200) {
                this.setState({flight: data.content})
            }
        });
    }

    render() {
        return (
            <VStack>
                {this.state.flight &&
                    <Box>
                        <Text>Рейс номер {this.flightNumber}</Text>
                        <Text>Аэропорт отправления: {this.state.flight?.fromAirport}</Text>
                        <Text>Аэропорт прибытия: {this.state.flight?.toAirport}</Text>
                        <Text>Время начала полёта: {this.state.flight?.date}</Text>
                        <Text>Стоимость: {this.state.flight?.price}</Text>
                    </Box>
                }
            </VStack>
        );
    };
};

export default FlightInfoPage;