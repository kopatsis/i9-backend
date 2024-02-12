package shared

import "math/rand"

func NameAnimals(reverese bool) string {
	animals := []string{
		"Aardvark", "Albatross", "Alligator", "Alpaca", "Ant", "Anteater", "Antelope",
		"Ape", "Armadillo", "Donkey", "Baboon", "Badger", "Barracuda", "Bat", "Bear",
		"Beaver", "Bee", "Bison", "Boar", "Buffalo", "Butterfly", "Camel", "Capybara",
		"Caribou", "Cassowary", "Cat", "Caterpillar", "Cattle", "Chameleon", "Cheetah",
		"Chicken", "Chimpanzee", "Chinchilla", "Cicada", "Clam", "Cobra", "Cockroach",
		"Cod", "Coyote", "Crab", "Crane", "Crocodile", "Crow", "Cuckoo", "Deer",
		"Dinosaur", "Dog", "Dolphin", "Duck", "Dugong", "Eagle", "Echidna", "Eel",
		"Elephant", "Elk", "Emu", "Falcon", "Ferret", "Finch", "Fish", "Flamingo", "Fly",
		"Fox", "Frog", "Gazelle", "Gerbil", "Giraffe", "Gnat", "Gnu", "Goat", "Goldfish",
		"Goose", "Gorilla", "Grasshopper", "Grouse", "Guanaco", "Gull", "Hamster",
		"Hare", "Hawk", "Hedgehog", "Heron", "Hippopotamus", "Hornet", "Horse", "Human",
		"Hummingbird", "Hyena", "Ibex", "Ibis", "Jackal", "Jaguar", "Jellyfish",
		"Kangaroo", "Kingfisher", "Koala", "Kookaburra", "Kouprey", "Kudu", "Lapwing",
		"Lark", "Lemur", "Leopard", "Lion", "Llama", "Lobster", "Locust", "Loris",
		"Louse", "Lyrebird", "Magpie", "Mallard", "Manatee", "Mandrill", "Mantis",
		"Marten", "Meerkat", "Mink", "Mole", "Mongoose", "Monkey", "Moose", "Mosquito",
		"Mouse", "Mule", "Narwhal", "Newt", "Nightingale", "Octopus", "Okapi", "Opossum",
		"Ostrich", "Otter", "Owl", "Ox", "Oyster", "Panda", "Panther", "Parrot",
		"Partridge", "Peafowl", "Pelican", "Penguin", "Pheasant", "Pig", "Pigeon",
		"Pony", "Porcupine", "Porpoise", "Quail", "Quelea", "Quetzal", "Rabbit",
		"Raccoon", "Rail", "Ram", "Rat", "Raven", "Reindeer",
		"Rhinoceros", "Rook", "Salamander", "Salmon", "Sand Dollar", "Sandpiper",
		"Sardine", "Scorpion", "Seahorse", "Seal", "Shark", "Sheep", "Shrew", "Skunk",
		"Snail", "Snake", "Sparrow", "Spider", "Spoonbill", "Squid", "Squirrel",
		"Starling", "Stingray", "Stinkbug", "Stork", "Swallow", "Swan", "Tapir",
		"Tarsier", "Termite", "Tiger", "Toad", "Trout", "Turkey", "Turtle", "Viper",
		"Vulture", "Wallaby", "Walrus", "Wasp", "Weasel", "Whale", "Wildcat", "Wolf",
		"Wolverine", "Wombat", "Woodcock", "Woodpecker", "Worm", "Wren", "Yak", "Zebra",
	}
	if reverese {
		w1 := animals[int(float32(len(animals)/2)*rand.Float32())]
		w2 := animals[len(animals)/2+int(float32(len(animals)/2)*rand.Float32())]
		return w1 + " " + w2
	}
	w1 := animals[len(animals)/2+int(float32(len(animals)/2)*rand.Float32())]
	w2 := animals[int(float32(len(animals)/2)*rand.Float32())]
	return w1 + " " + w2
}
