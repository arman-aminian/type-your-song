package store

import (
	"context"
	"github.com/arman-aminian/type-your-song/model"
	"github.com/arman-aminian/type-your-song/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore struct {
	db *mongo.Collection
}

func NewUserStore(db *mongo.Collection) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) Create(u *model.User) error {
	_, err := us.db.InsertOne(context.TODO(), u)
	return err
}

func (us *UserStore) Remove(field, value string) error {
	_, err := us.db.DeleteOne(context.TODO(), bson.M{field: value})
	return err
}

func (us *UserStore) UpdateStrField(old *model.User, field string, value string) error {
	var err error
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": old.ID}, bson.M{"$set": bson.M{field: value}})
	return err
}

func (us *UserStore) UpdateStrFieldByEmail(old *model.User, field string, value string) error {
	var err error
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"email": old.Email}, bson.M{"$set": bson.M{field: value}})
	return err
}

func (us *UserStore) UpdateBoolField(old *model.User, field string, value bool) error {
	var err error
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": old.ID}, bson.M{"$set": bson.M{field: value}})
	return err
}

func (us *UserStore) UpdateBoolFieldByEmail(old *model.User, field string, value bool) error {
	var err error
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"email": old.Email}, bson.M{"$set": bson.M{field: value}})
	return err
}

func (us *UserStore) UpdateProfile(u *model.User) error {
	_, err := us.db.UpdateOne(context.TODO(),
		bson.M{"_id": u.Username},
		bson.M{"$set": bson.M{
			//"name": u.Name,
			//"bio":             u.Bio,
			//"profile_picture": u.ProfilePicture,
			//"header_picture":  u.HeaderPicture,
		},
		})
	return err
}

func (us *UserStore) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (us *UserStore) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"username": username}).Decode(&u)
	return &u, err
}

func (us *UserStore) GetById(id primitive.ObjectID) (*model.User, error) {
	var u model.User
	err := us.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&u)
	return &u, err
}

func (us *UserStore) AddFollowing(current primitive.ObjectID, u primitive.ObjectID) (model.User, error) {
	cu, err := us.GetById(current)
	if err != nil {
		return model.User{}, err
	}
	*cu.Followings = append(*cu.Followings, u)
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": current}, bson.M{"$set": bson.M{"followings": cu.Followings}})
	return *cu, err
}

func (us *UserStore) RemoveFollowing(current primitive.ObjectID, u primitive.ObjectID) (model.User, error) {
	cu, err := us.GetById(current)
	if err != nil {
		return model.User{}, err
	}
	newFollowings := &[]primitive.ObjectID{}
	for _, o := range *cu.Followings {
		if o != u {
			*newFollowings = append(*newFollowings, o)
		}
	}
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": current}, bson.M{"$set": bson.M{"followings": newFollowings}})
	if err != nil {
		return model.User{}, err
	}
	cu.Followings = newFollowings
	return *cu, nil
}

func (us *UserStore) addScore(uid primitive.ObjectID, score int) error {
	u, err := us.GetById(uid)
	if err != nil {
		return err
	}
	u.Score = u.Score + score
	_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"score": u.Score}})
	return err
}

func (us *UserStore) Record(uid primitive.ObjectID, passed model.PassedSong, s *model.Song) (int, error) {
	u, err := us.GetById(uid)
	if err != nil {
		return 0, err
	}

	score := 0
	ps := *u.PassedSongs
	checked := false
	for i, p := range ps {
		if p.SID == passed.SID {
			checked = true
			if utils.LevelToNum(p.PassedLevel) < utils.LevelToNum(passed.PassedLevel) {
				score = calculateCurrentScore(s, passed) - calculateCurrentScore(s, p)
				ps = append(ps[:i], ps[i+1:]...)
				ps = append(ps, passed)
				_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"passed_songs": &ps}})
				break
			} else if utils.LevelToNum(p.PassedLevel) == utils.LevelToNum(passed.PassedLevel) {
				if p.Speed < passed.Speed {
					score = calculateCurrentScore(s, passed) - calculateCurrentScore(s, p)
					ps = append(ps[:i], ps[i+1:]...)
					ps = append(ps, passed)
					_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"passed_songs": &ps}})
					break
				} else if p.Speed == passed.Speed {
					if p.Accuracy < passed.Accuracy {
						score = calculateCurrentScore(s, passed) - calculateCurrentScore(s, p)
						ps = append(ps[:i], ps[i+1:]...)
						ps = append(ps, passed)
						_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"passed_songs": &ps}})
						break
					}
				}
			}
		}
	}
	if !checked {
		score = calculateCurrentScore(s, passed)
		*u.PassedSongs = append(*u.PassedSongs, passed)
		_, err = us.db.UpdateOne(context.TODO(), bson.M{"_id": uid}, bson.M{"$set": bson.M{"passed_songs": u.PassedSongs}})
	}
	if score < 0 {
		score = 0
	}
	return score, err
}
func calculateCurrentScore(s *model.Song, p model.PassedSong) int {
	score := utils.CalculateScore(s.WordsCount, s.MaxWPM, s.AvgWPM, utils.LevelToNum(p.PassedLevel))
	if p.PassedLevel == utils.Simple {
		score = score * p.Speed * p.Accuracy / 100 / s.AvgWPM
	}
	return score
}
