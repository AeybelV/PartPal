import argparse
import npyscreen
import concurrent.futures

import partpal.core.config as pp_config
from partpal.core.bom_parser import parse_bom_csv

# TODO: Use the Distributor interface instead of specific ones
from partpal.distributors.digikey import DigiKeyDistributor
from partpal.distributors.mouser import MouserDistributor


def parse_args():
    """Parse the CLI args"""
    parser = argparse.ArgumentParser(description="PoorPartPal")
    parser.add_argument("-i", "--input", type=str, help="Filepath to input BOM")
    parser.add_argument(
        "-o",
        "--optimize",
        help="Optimize the BOM",
    )

    args = parser.parse_args()
    return args


class BOMGrid(npyscreen.SimpleGrid):
    def custom_print_cell(self, cell, value):
        value_str = str(value) if value else ""
        max_width = cell.width - 1  # Leave space for padding or borders
        cell.value = value_str[:max_width] if len(value_str) > max_width else value_str


class BOMView(npyscreen.FormBaseNew):
    def create(self):
        # Get the maximum usable width and height for the form
        screen_max_y, screen_max_x = self.useable_space()

        self.name = "PoorPartPal"
        self.bom_data = self.parentApp.bom_data

        # Top pane: Buttons and actions
        self.actions_pane = self.add(
            npyscreen.BoxTitle, name="Actions", max_height=int(0.1 * screen_max_y)
        )
        self.actions_pane.values = ["Optimize BOM", "Exit"]
        self.actions_pane.when_value_edited = self.handle_action_selection

        # Add column titles manually as a separate line
        self.add(
            npyscreen.FixedText,
            value="Name              Quantity      Part Number     Description     Cost          Distributor",
            editable=False,
            relx=1,
            rely=int(0.15 * screen_max_y),
            width=screen_max_x - 2,
        )

        # Middle pane: Grid to display the BOM data
        self.table_pane = self.add(
            BOMGrid,
            relx=1,
            rely=int(0.17 * screen_max_y),
            max_height=int(0.75 * screen_max_y),
            column_width=int(0.1 * screen_max_x),
        )
        self.table_pane.values = self.bom_to_grid_values(self.bom_data)

        self.total_cost = self.calculate_total_cost()

        # Bottom pane: Keyboard shortcut hints
        self.hint_pane = self.add(
            npyscreen.FixedText,
            value=f"Shortcuts: Ctrl+T = Table, Ctrl+A = Actions, Ctrl+Q = Quit | Total Cost: {self.total_cost:.4f} USD",
            editable=False,
            rely=int(0.95 * screen_max_y),
            relx=2,
        )

        # Add custom key bindings
        self.add_handlers(
            {
                "^T": self.switch_to_table,  # Ctrl+T to switch to the table
                "^A": self.switch_to_actions,  # Ctrl+A to switch to the actions pane
                "^Q": self.exit_application,  # Ctrl+Q to exit
            }
        )

    def calculate_total_cost(self):
        # Calculate the total cost from the BOM data
        total_cost = 0.0
        for component in self.bom_data:
            try:
                total_cost += float(component.get("cost", 0.0)) * int(
                    component.get("quantity", 1)
                )
            except ValueError:
                continue  # Skip if there's an issue with the data format
        return total_cost

    def switch_to_table(self, *args):
        self.set_editing(self.table_pane)
        self.table_pane.edit()

    def switch_to_actions(self, *args):
        self.set_editing(self.actions_pane)
        self.actions_pane.edit()

    def bom_to_grid_values(self, bom_data):
        return [
            [
                component.get("name", ""),
                component.get("quantity", ""),
                component.get("part_number", ""),
                component.get("description", ""),
                component.get("cost", ""),
                component.get("distributor", ""),
            ]
            for component in bom_data
        ]

    def handle_action_selection(self):
        if self.actions_pane.value is not None:
            action = self.actions_pane.values[self.actions_pane.value]
            if action == "Optimize BOM":
                self.optimize_bom()
            elif action == "Exit":
                self.exit_application()
        else:
            pass

    def update_component(self, component, distributor):
        status, response = distributor.get_product_information(component["part_number"])

        if status == True:
            cost = response["unitPrice"]
            current_cost = component["cost"]

            try:
                current_cost = float(current_cost)
            except:
                current_cost = ""

            if current_cost == "" or cost < current_cost:
                component["distributor"] = distributor.name
                component["cost"] = str(cost)

        return component

    def optimize_bom(self):

        # TODO: Change to be general Distributor, not just digikey and mouser
        config = pp_config.load_config()

        digikey = DigiKeyDistributor(
            config["distributors"]["digikey"]["client_id"],
            config["distributors"]["digikey"]["client_secret"],
            config["distributors"]["digikey"]["sandbox"],
        )
        auth_status = digikey.get_access_token()

        mouser = MouserDistributor(config["distributors"]["mouser"]["api_key"])

        dists = [digikey, mouser]

        for distributor in dists:
            # Optimization logic can go here

            # Creates a threadpool, each component get its own thread to query the distributor
            # Meaning API requests are done in parallel, they will update the component and bom data
            with concurrent.futures.ThreadPoolExecutor() as executor:
                # TODO: Change to be general Distributor, not just digikey
                self.bom_data = list(
                    executor.map(
                        lambda c: self.update_component(c, distributor), self.bom_data
                    )
                )

            # After the optimization, refresh the table view
            self.table_pane.values = self.bom_to_grid_values(self.parentApp.bom_data)
            self.table_pane.update()

            # Recalculate the total cost after optimization
            self.total_cost = self.calculate_total_cost()

            # Update the hint pane with the new total cost
            self.hint_pane.value = f"Shortcuts: Ctrl+T = Table, Ctrl+A = Actions, Ctrl+Q = Quit | Total Cost: {self.total_cost:.4f} USD"
            self.hint_pane.update()

        npyscreen.notify_confirm("Optimization complete!", title="Optimize")

        # Reset the actions pane selection to allow re-execution
        self.actions_pane.entry_widget.value = None
        self.actions_pane.display()

    def exit_application(self, *args, **kwargs):
        self.parentApp.setNextForm(None)
        self.parentApp.switchFormNow()


class PartPalTUI(npyscreen.NPSAppManaged):
    def __init__(self, bom_file=None, optimize=False, **kwargs):
        super().__init__(**kwargs)
        self.bom_file = bom_file
        self.optimize = optimize
        self.bom_data = None

    def onStart(self):
        if self.bom_file:
            self.bom_data = parse_bom_csv(self.bom_file)
            self.show_bom_view()
        else:
            self.show_splash_screen()

    def show_splash_screen(self):
        self.addForm("MAIN", SplashScreen, name="PartPal TUI")

    def show_bom_view(self):
        self.addForm(
            "MAIN",
            BOMView,
            name="BOM Viewer",
            bom_data=self.bom_data,
            optimize=self.optimize,
        )


class SplashScreen(npyscreen.Form):
    def create(self):
        self.add(
            npyscreen.TitleText,
            name="Welcome to PartPal!",
            value="Press Enter to continue.",
        )
        self.add_handlers({"^Q": self.exit_application})

    def afterEditing(self):
        self.parentApp.setNextForm(None)

    def exit_application(self, *args, **kwargs):
        self.parentApp.setNextForm(None)
        self.parentApp.switchFormNow()


def run_tui():
    args = parse_args()

    if args.input:
        # Start the app with a BOM file
        PartPalTUI(bom_file=args.input, optimize=args.optimize).run()
    else:
        # Start the app without a BOM file
        PartPalTUI().run()


if __name__ == "__main__":
    run_tui()
